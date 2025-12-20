# 2Web Kit - Shared

Shared utilities and types for 2Web packages.

These functions are not intended for consumption, and should not be exported.

The benefit of using shared code is so that multiple packages can use the same
global store for things like DOM updates, meaning that each package can use the
same update queue rather than having separate queues that may conflict with each
other.
