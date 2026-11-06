from fastapi.testclient import TestClient

from internal.adapters.inbound.http import create_app
from internal.adapters.outbound.sql import (
    SqlAlchemyBoardRepository,
    SqlAlchemyTaskRepository,
    create_db_engine,
    create_session_factory,
    init_db,
)
from internal.core.application.usecases.board import (
    AddColumnService,
    CreateBoardService,
    GetBoardService,
    ListBoardsService,
)
from internal.core.application.usecases.task import (
    CreateTaskService,
    DeleteTaskService,
    ListTasksService,
    MoveTaskService,
)


def create_test_client():
    engine = create_db_engine("sqlite:///:memory:")
    init_db(engine)
    session_factory = create_session_factory(engine)

    board_repo = SqlAlchemyBoardRepository(session_factory)
    task_repo = SqlAlchemyTaskRepository(session_factory)

    app = create_app(
        create_board_uc=CreateBoardService(board_repo),
        get_board_uc=GetBoardService(board_repo),
        list_boards_uc=ListBoardsService(board_repo),
        add_column_uc=AddColumnService(board_repo),
        create_task_uc=CreateTaskService(board_repo, task_repo),
        move_task_uc=MoveTaskService(board_repo, task_repo),
        list_tasks_uc=ListTasksService(board_repo, task_repo),
        delete_task_uc=DeleteTaskService(task_repo),
    )
    return TestClient(app)


def test_http_board_and_task_lifecycle():
    client = create_test_client()

    # 1. POST /api/v1/boards
    resp = client.post("/api/v1/boards", json={"name": "Sprint 1", "description": "Release MVP"})
    assert resp.status_code == 201
    board_data = resp.json()
    board_id = board_data["id"]
    assert board_data["name"] == "Sprint 1"
    assert len(board_data["columns"]) == 3
    in_progress_col = [c for c in board_data["columns"] if c["name"] == "In Progress"][0]

    # 2. GET /api/v1/boards
    resp = client.get("/api/v1/boards")
    assert resp.status_code == 200
    assert len(resp.json()) == 1

    # 3. GET /api/v1/boards/{boardId}
    resp = client.get(f"/api/v1/boards/{board_id}")
    assert resp.status_code == 200
    assert resp.json()["id"] == board_id

    # 4. POST /api/v1/boards/{boardId}/columns
    resp = client.post(
        f"/api/v1/boards/{board_id}/columns",
        json={"name": "QA Review", "wip_limit": 2},
    )
    assert resp.status_code == 201
    qa_col = resp.json()
    qa_col_id = qa_col["id"]
    assert qa_col["name"] == "QA Review"
    assert qa_col["wip_limit"] == 2

    # 5. POST /api/v1/boards/{boardId}/tasks
    resp1 = client.post(
        f"/api/v1/boards/{board_id}/tasks",
        json={"column": "Backlog", "title": "Task 1", "priority": "HIGH"},
    )
    assert resp1.status_code == 201
    t1_id = resp1.json()["id"]

    resp2 = client.post(
        f"/api/v1/boards/{board_id}/tasks",
        json={"column": "Backlog", "title": "Task 2", "priority": "MEDIUM"},
    )
    assert resp2.status_code == 201
    t2_id = resp2.json()["id"]

    resp3 = client.post(
        f"/api/v1/boards/{board_id}/tasks",
        json={"column": "Backlog", "title": "Task 3", "priority": "LOW"},
    )
    assert resp3.status_code == 201
    t3_id = resp3.json()["id"]

    # 6. GET /api/v1/boards/{boardId}/tasks
    resp = client.get(f"/api/v1/boards/{board_id}/tasks")
    assert resp.status_code == 200
    assert len(resp.json()) == 3

    # 7. PATCH /api/v1/tasks/{taskId}/move - Move task 1 and 2 to QA Review (WIP limit 2)
    resp = client.patch(f"/api/v1/tasks/{t1_id}/move", json={"target_column_id": qa_col_id})
    assert resp.status_code == 200
    assert resp.json()["column_name"] == "QA Review"

    resp = client.patch(f"/api/v1/tasks/{t2_id}/move", json={"target_column_id": qa_col_id})
    assert resp.status_code == 200

    # 8. Move task 3 to QA Review -> Should return 409 Conflict with WIP_LIMIT_EXCEEDED
    resp = client.patch(f"/api/v1/tasks/{t3_id}/move", json={"target_column_id": qa_col_id})
    assert resp.status_code == 409
    err = resp.json()["error"]
    assert err["code"] == "WIP_LIMIT_EXCEEDED"
    assert err["current_count"] == 2
    assert err["limit"] == 2
    assert err["column_id"] == qa_col_id
    assert "Cannot move task to 'QA Review'" in err["message"]

    # 9. DELETE /api/v1/tasks/{taskId}
    resp = client.delete(f"/api/v1/tasks/{t1_id}")
    assert resp.status_code == 204

    # Now QA Review has 1 task; task 3 can move!
    resp = client.patch(f"/api/v1/tasks/{t3_id}/move", json={"target_column_id": qa_col_id})
    assert resp.status_code == 200

    # 10. 404 Tests
    resp = client.get("/api/v1/boards/non-existent-id")
    assert resp.status_code == 404
    assert resp.json()["error"]["code"] == "BOARD_NOT_FOUND"

    resp = client.delete(f"/api/v1/tasks/{t1_id}")
    assert resp.status_code == 404
    assert resp.json()["error"]["code"] == "TASK_NOT_FOUND"

