import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { QueryClient } from "@tanstack/react-query";
import React, { Suspense } from "react";
import { ThemeProvider } from "@/provider/theme-provider.tsx";
import { ModeToggle } from "@/components/mode-toggle.tsx";

export const Route = createRootRouteWithContext<{
  queryClient: QueryClient;
}>()({
  component: RootComponent,
});

const TanStackRouterDevtools =
  process.env.NODE_ENV === "production"
    ? () => null // Render nothing in production
    : React.lazy(() =>
        // Lazy load in development
        import("@tanstack/router-devtools").then((res) => ({
          default: res.TanStackRouterDevtools,
          // For Embedded Mode
          // default: res.TanStackRouterDevtoolsPanel
        })),
      );

function RootComponent() {
  return (
    <ThemeProvider defaultTheme="system" storageKey="gitlab-classrooms-theme">
      <div className="h-screen">
        <ModeToggle />
        <div className="max-w-2xl m-auto">
          <Outlet />
          <ReactQueryDevtools initialIsOpen={false} />
          <Suspense>
            <TanStackRouterDevtools />
          </Suspense>
        </div>
      </div>
    </ThemeProvider>
  );
}
