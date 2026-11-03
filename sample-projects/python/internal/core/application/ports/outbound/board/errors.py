class OutboundRepositoryError(Exception):
    """Raised when an unexpected outbound storage or infrastructure failure occurs."""
    def __init__(self, message: str):
        super().__init__(message)
        self.message = message

