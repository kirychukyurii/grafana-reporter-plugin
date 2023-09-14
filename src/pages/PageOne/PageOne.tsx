import * as React from "react";
import { getBackendSrv } from "@grafana/runtime";
import { testIds } from "../../components/testIds";
import { useAsync } from "react-use";
import { Badge, HorizontalGroup } from "@grafana/ui";

export const PageOne = () => {
    const { error, loading, value } = useAsync(() => {
        const backendSrv = getBackendSrv();

        return Promise.all([
            backendSrv.get(`api/plugins/kirychukyurii-reporter-app/resources/ping`),
            backendSrv.get(`api/plugins/kirychukyurii-reporter-app/health`),
        ]);
    });

    if (loading) {
        return (
            <div data-testid={testIds.schedule.container}>
                <span>Loading...</span>
            </div>
        );
    }

    if (error || !value) {
        return (
            <div data-testid={testIds.schedule.container}>
                <span>Error: {error?.message}</span>
            </div>
        );
    }

    const [ping, health] = value;

    return (
        <div data-testid={testIds.schedule.container}>
            <HorizontalGroup>
                <h3>Plugin Health Check</h3>{" "}
                <span data-testid={testIds.schedule.health}>
          {renderHealth(health?.message)}
        </span>
            </HorizontalGroup>
            <HorizontalGroup>
                <h3>Ping Backend</h3>{" "}
                <span data-testid={testIds.schedule.ping }>{ping?.message}</span>
            </HorizontalGroup>
        </div>
    );
};

function renderHealth(message: string | undefined) {
    switch (message) {
        case "ok":
            return <Badge color="green" text="OK" icon="heart" />;

        default:
            return <Badge color="red" text="BAD" icon="bug" />;
    }
}
