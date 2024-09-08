import { Route, Routes, useLocation } from "react-router-dom"
import { CSSTransition, TransitionGroup } from 'react-transition-group';
import { Auth } from "../../pages/auth";
import { ErrorPage } from "../../pages/errorPage";

export const Router = () => {
    const location = useLocation();

    return (
        <TransitionGroup>
            <CSSTransition
                key={location.pathname}
                classNames="fade"
                timeout={300}
            >
                <Routes location={location}>
                    <Route index path="/"/>
                    <Route path="/auth" element={<Auth />} />

                    <Route path="/app">
                        <Route path="/app" />
                        <Route path="/app/tasks" />
                        <Route path="/app/time-blocks" />
                        <Route path="/app/life-calendar" />
                    </Route>

                    <Route path="*" element={<ErrorPage status={404} message="Страница не найдена!" />} />
                </Routes>
            </CSSTransition>
        </TransitionGroup>
    )
}