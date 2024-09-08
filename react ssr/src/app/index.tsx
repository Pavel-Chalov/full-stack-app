import { Providers } from "./providers"
import { Router } from "./router"
import "./styles/index.scss"

export const App = () => {
    return (
        <Providers>
            <Router />
        </Providers>
    )
}