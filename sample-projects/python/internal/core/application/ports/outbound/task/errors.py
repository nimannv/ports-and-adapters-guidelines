class OutboundTaskRepositoryError(Exception):
    """Raised when an unexpected outbound storage error occurs for task persistence."""
    def __init__(self, message: str):
        super().__init__(message)
        self.message = message

