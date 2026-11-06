from datetime import datetime, timezone
from sqlalchemy import Column, DateTime, ForeignKey, Integer, String, Text
from sqlalchemy.orm import declarative_base, relationship

Base = declarative_base()


class BoardModel(Base):
    __tablename__ = "boards"

    id = Column(String(36), primary_key=True)
    name = Column(String(255), nullable=False)
    description = Column(Text, default="", nullable=False)
    created_at = Column(DateTime(timezone=True), default=lambda: datetime.now(timezone.utc), nullable=False)
    updated_at = Column(DateTime(timezone=True), default=lambda: datetime.now(timezone.utc), nullable=False)

    columns = relationship(
        "ColumnModel",
        back_populates="board",
        cascade="all, delete-orphan",
        order_by="ColumnModel.order_index",
        lazy="joined",
    )
    tasks = relationship(
        "TaskModel",
        back_populates="board",
        cascade="all, delete-orphan",
    )


class ColumnModel(Base):
    __tablename__ = "board_columns"

    id = Column(String(36), primary_key=True)
    board_id = Column(String(36), ForeignKey("boards.id", ondelete="CASCADE"), nullable=False, index=True)
    name = Column(String(255), nullable=False)
    order_index = Column(Integer, nullable=False, default=0)
    wip_limit = Column(Integer, nullable=True)

    board = relationship("BoardModel", back_populates="columns")


class TaskModel(Base):
    __tablename__ = "tasks"

    id = Column(String(36), primary_key=True)
    board_id = Column(String(36), ForeignKey("boards.id", ondelete="CASCADE"), nullable=False, index=True)
    column_id = Column(String(36), ForeignKey("board_columns.id", ondelete="CASCADE"), nullable=False, index=True)
    title = Column(String(255), nullable=False)
    description = Column(Text, default="", nullable=False)
    priority = Column(String(50), nullable=False, default="MEDIUM")
    due_date = Column(DateTime(timezone=True), nullable=True)
    created_at = Column(DateTime(timezone=True), default=lambda: datetime.now(timezone.utc), nullable=False)
    updated_at = Column(DateTime(timezone=True), default=lambda: datetime.now(timezone.utc), nullable=False)

    board = relationship("BoardModel", back_populates="tasks")

