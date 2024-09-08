import { FC, memo, useState } from "react";
import "./authForms.scss"
import { animateForm } from "../lib/animateForm";
import { Title } from "../../../shared/ui/title";
import { Input } from "../../../shared/ui/input";
import { Button } from "../../../shared/ui/button";
import { toast } from "react-toastify";
import { useMutation } from "@tanstack/react-query";
import { signInRequest } from "../../../entities/user/api/signInRequest";
import { AxiosError, AxiosResponse } from "axios";
import { WithMessage } from "../../../shared/types/withMessage";
import { useNavigate } from "react-router-dom";

export const SignInForm: FC = memo(() => {
    const [name, setName] = useState("")
    const [password, setPassword] = useState("")
    const [nameError, setNameError] = useState<string | null>(null)
    const [passwordError, setPasswordError] = useState<string | null>(null)

    const navigate = useNavigate()

    function handleOnChange(e: React.ChangeEvent<HTMLInputElement>) {
        const value = e.target.value.replace(" ", "").substring(0, 48)

        switch (e.target.name) {
            case "name":
                setName(value.substring(0, 24))
                setNameError(null)
                break;
            case "password":
                setPassword(value)
                setPasswordError(null)
                break
        }
    }

    const {mutate: signIn} = useMutation<AxiosResponse<WithMessage>, AxiosError<WithMessage>>({
        mutationKey: ["signIn"],
        mutationFn: () => signInRequest({name: name, password: password}),
        onSuccess: (res) => {
            toast.success(res.data.message)
            navigate("/app")
        },
        onError: (err) => {
            if (err.response!.status === 404) {
                setNameError(err.response!.data.message)
            } else if (err.response!.status === 403) {
                console.log(err.response!.data.message)
                setPasswordError(err.response!.data.message)
            } else if (err.response === undefined) {
                toast.error("Сервер не работает. Попробуйте позже!")
            }
        }
    })

    function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault()
        
        if (name.length >= 3 && password.length >= 8 && !nameError && !passwordError) {
            signIn()
        }
    }

    return (
        <form onSubmit={handleSubmit} style={{display: "none"}} className="auth-form">
            <Title Level="h1" className="auth-form__title">Вход</Title>

            <Input error={nameError!} uiSize="big" value={name} onChange={handleOnChange} name="name" className="auth-form__field" label="Логин" type="text" placeholder="Введите логин" />
            <Input error={passwordError!} uiSize="big" value={password} onChange={handleOnChange} name="password" className="auth-form__field" label="Пароль" type="password" placeholder="Введите пароль" />
            <Button uiSize="big" className="auth-form__button" disabled={nameError || passwordError || name.length < 3 || password.length < 8 ? true : false} type="submit">Войти</Button>

            <p>Еще не зарегистрированы? <span onClick={animateForm}>Создать аккаунт</span></p>
        </form>
    )
})