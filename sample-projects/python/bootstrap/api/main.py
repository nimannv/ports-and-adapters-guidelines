import sys
import uvicorn

from configs.config import load_config
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


def build_app():
    # 1. Configuration
    config = load_config()

    # 2. Outbound Adapters (Database & Repositories)
    engine = create_db_engine(config.database_url)
    init_db(engine)
    session_factory = create_session_factory(engine)

    board_repo = SqlAlchemyBoardRepository(session_factory)
    task_repo = SqlAlchemyTaskRepository(session_factory)

    # 3. Application Use Cases
    create_board_uc = CreateBoardService(board_repo)
    get_board_uc = GetBoardService(board_repo)
    list_boards_uc = ListBoardsService(board_repo)
    add_column_uc = AddColumnService(board_repo)

    create_task_uc = CreateTaskService(board_repo, task_repo)
    move_task_uc = MoveTaskService(board_repo, task_repo)
    list_tasks_uc = ListTasksService(board_repo, task_repo)
    delete_task_uc = DeleteTaskService(task_repo)

    # 4. Inbound Adapter (HTTP / FastAPI)
    app = create_app(
        create_board_uc=create_board_uc,
        get_board_uc=get_board_uc,
        list_boards_uc=list_boards_uc,
        add_column_uc=add_column_uc,
        create_task_uc=create_task_uc,
        move_task_uc=move_task_uc,
        list_tasks_uc=list_tasks_uc,
        delete_task_uc=delete_task_uc,
    )
    return app, config


app, app_config = build_app()


def run():
    uvicorn.run(
        "bootstrap.api.main:app",
        host=app_config.api_host,
        port=app_config.api_port,
        log_level=app_config.log_level.lower(),
        reload=False,
    )


if __name__ == "__main__":
    run()

