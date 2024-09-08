import { FC } from "react"
import { Title } from "../../../shared/ui/title"
import "./errorPage.scss"
import { Button } from "../../../shared/ui/button"
import { useNavigate } from "react-router-dom"

interface ErrorPageProps {
    status: number,
    message: string,
}

export const ErrorPage: FC<ErrorPageProps> = ({status, message}) => {
    const navigate = useNavigate()

    return (
        <div className="error-page page centered">
            <Title className="error-page__title" Level="h1">{status}</Title>
            <Title className="error-page__sub-title" Level="h2">{message}</Title>
            <Button onClick={() => navigate("/")} uiSize="big">На главную</Button>
        </div>
    )
}