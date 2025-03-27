#!/usr/bin/env python3
# Copyright 2025 familia
# See LICENSE file for licensing details.

"""FastAPI Charm entrypoint."""

import logging
import os
import typing

import ops
import paas_charm.fastapi
from ops.model import ActiveStatus, WaitingStatus

logger = logging.getLogger(__name__)


class OpheliaCiInterfaceCharm(paas_charm.fastapi.Charm):
    """FastAPI Charm service."""

    _server_address = None

    def __init__(self, *args: typing.Any) -> None:
        """Initialize the instance.

        Args:
            args: passthrough to CharmBase.
        """
        super().__init__(*args)
        self.framework.observe(
            self.on.ophelia_server_relation_joined, self._on_ophelia_server_joined
        )
        self.framework.observe(
            self.on.ophelia_server_relation_changed, self._on_ophelia_server_changed
        )
        self.framework.observe(self.on.start, self._on_start)

    def _on_start(self, event: ops.framework.EventBase):
        """Handle the start event."""
        if not self._server_address:
            self.unit.status = WaitingStatus("Waiting for ophelia-server address")
            return
        self.unit.status = ActiveStatus(f"Ready to connect to {self._server_address}")
        self._publish_client_info()

    def _on_ophelia_server_relation_joined(self, event: ops.framework.EventBase):
        """Handle the ophelia-server relation joined event."""
        self.unit.status = WaitingStatus("Waiting for ophelia-server address")

    def _on_ophelia_server_relation_changed(self, event: ops.framework.EventBase):
        """Handle the ophelia-server relation changed event."""
        if not event.relation.data.get(event.app):
            logger.warning("No application data found in ophelia-server relation.")
            self.unit.status = WaitingStatus("Waiting for ophelia-server address")
            return

        server_address = event.relation.data[event.app].get("server_address")
        server_version = event.relation.data[event.app].get("server_version")

        if not server_address:
            logger.warning("Server address not found in relation data.")
            self.unit.status = WaitingStatus("Waiting for ophelia-server address")

        logger.info(f"Received server address: {server_address} (version: {server_version})")
        self._server_address = server_address
        os.environ["APP_OPHELIA_CI_GRPC_SERVER"] = self._server_address
        self.unit.status = ActiveStatus(f"Ready to connect to {self._server_address}")
        self._publish_client_info(event.relation)

    def _publish_client_info(self, relation: typing.Optional[ops.Relation] = None):
        """Publish the client's address and version on the relation."""
        if not relation:
            for rel in self.model.relations["ophelia-server"]:
                self._publish_client_info(rel)
            return

        client_address = (
            f"{self.unit.name.split('/')[0]}:8008"
        )
        client_version = "1"

        relation.data[self.unit]["client_address"] = client_address
        relation.data[self.unit]["client_version"] = client_version
        logger.info(f"Published client info: address={client_address}, version={client_version}")


if __name__ == "__main__":
    ops.main(OpheliaCiInterfaceCharm)
