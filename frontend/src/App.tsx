import { lazy, Suspense } from "solid-js";
import { Router, Route } from "@solidjs/router";

import LandingPage from "./pages/LandingPage";
const AppPage = lazy(() => import("./pages/AppPage"));

function App() {
  return (
    <Router>
        <Route path="/" component={LandingPage} />
        <Route path="/app" component={() => (
          <Suspense fallback={<div class="min-h-screen bg-black text-green-400 font-mono flex items-center justify-center">Loading...</div>}>
            <AppPage />
          </Suspense>
        )} />
    </Router>
  );
}

export default App;