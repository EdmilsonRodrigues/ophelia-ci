#!/usr/bin/env python3
# Copyleft 2025 familia
# See LICENSE file for licensing details.

"""Go Charm entrypoint."""

import logging
import os
import typing

import ops
import paas_charm.go
from ops.model import ActiveStatus

logger = logging.getLogger(__name__)


class OpheliaCiServerCharm(paas_charm.go.Charm):
    """Go Charm service."""

    def __init__(self, *args: typing.Any) -> None:
        """Initialize the instance.

        Args:
            args: passthrough to CharmBase.
        """
        super().__init__(*args)
        self.git_repos_storage = ops.Storage("git-repos")
        self.framework.observe(self.git_repos_storage.on.attached, self._on_git_repos_attached)
        self.framework.observe(
            self.on.ophelia_interface_relation_joined, self._on_ophelia_interface_joined
        )

    def _on_git_repos_attached(self, event: ops.framework.EventBase):
        """Handle the git-repos storage being attached."""
        if not event.storage:
            logger.warning("Git repos storage attached event received without storage object.")
            return

        mount_point = event.storage.location
        logger.info(f"Git repos storage attached at: {mount_point}")
        self.unit.status = ActiveStatus(f"Git repos storage attached at {mount_point}")
        os.environ["APP_OPHELIA_CI_SERVER_HOME_PATH"] = mount_point

    def _on_ophelia_interface_relation_joined(self, event: ops.framework.EventBase):
        """Handle the ophelia-interface relation being established."""
        self._publish_server_info(event.relation)

    def _publish_server_info(self, relation: ops.Relation):
        """Publish the server's address and version on the relation."""
        server_address = f"{self.unit.name.split('/')[0]}:50051"  # Example: unit-id:port
        server_version = "1"

        relation.data[self.unit]["server_address"] = server_address
        relation.data[self.unit]["server_version"] = server_version
        self.unit.status = ActiveStatus(
            f"Providing server info: {server_address} (v{server_version})"
        )


if __name__ == "__main__":
    ops.main(OpheliaCiServerCharm)
