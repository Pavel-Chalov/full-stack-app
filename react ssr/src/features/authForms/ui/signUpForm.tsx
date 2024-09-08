import { FC, useState } from "react";
import "./authForms.scss"
import { animateForm } from "../lib/animateForm";
import { Title } from "../../../shared/ui/title";
import { Input } from "../../../shared/ui/input";
import { Button } from "../../../shared/ui/button";
import { useMutation } from "@tanstack/react-query";
import { signUpRequest } from "../../../entities/user/api/signUpRequest";
import { AxiosError, AxiosResponse } from "axios";
import { WithMessage } from "../../../shared/types/withMessage";
import { toast } from "react-toastify";
import { useNavigate } from "react-router-dom";

export const SignUpForm: FC = () => {
    const [name, setName] = useState("")
    const [password, setPassword] = useState("")
    const [repeat, setRepeat] = useState("")

    const [nameError, setNameError] = useState<null | string>("Имя должно быть не менее 3 символов!")
    const [passwordError, setPasswordError] = useState<null | string>("Пароль должен быть не менее 8 символов!")
    const [repeatError, setRepeatError] = useState<null | string>(null)

    const [showErrors, setShowErrors] = useState(false)

    const navigate = useNavigate()

    function handleOnChange(e: React.ChangeEvent<HTMLInputElement>) {
        const value = e.target.value.replace(" ", "").substring(0, 48)

        switch (e.target.name) {
            case "name":
                setName(value.substring(0, 24))

                if (value.length < 3) {
                    setNameError("Имя должно быть не менее 3 символов!")
                } else {
                    setNameError(null)
                }
                
                break
            case "password":
                setPassword(value)

                if (value.length < 8) {
                    setPasswordError("Пароль должен быть не менее 8 символов!")
                } else {
                    setPasswordError(null)
                }

                if (value !== repeat) {
                    setRepeatError("Пароли не совпадают!")
                } else {
                    setRepeatError(null)
                }

                break
            case "repeat":
                setRepeat(value)

                if (value !== password) {
                    setRepeatError("Пароли не совпадают!")
                } else {
                    setRepeatError(null)
                }

                break
        }
    }

    const {mutate: signUp} = useMutation<AxiosResponse<WithMessage>, AxiosError<WithMessage>>({
        mutationKey: ["signUp"],
        mutationFn: () => signUpRequest({name: name,  password: password}),
        onSuccess: (res) => {
            toast.success(res.data.message)
            navigate("/app")
        },
        onError: (err) => {
            if (err.response !== undefined) {
                setNameError(err.response!.data.message)
            } else if (err.response === undefined) {
                toast.error("Сервер не работает. Попробуйте позже!")
            }
        }
    })

    function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault()
        setShowErrors(true)

        if (!nameError && !passwordError && !repeatError) {
            signUp()
        }
    }

    return (
        <form onSubmit={handleSubmit} className="auth-form">
            <Title Level="h1" className="auth-form__title">Регистрация</Title>

            <Input uiSize="big" error={showErrors ? nameError! : undefined} value={name} onChange={handleOnChange} name="name" className="auth-form__field" label="Логин" type="text" placeholder="Введите логин" />
            <Input uiSize="big" error={showErrors ? passwordError! : undefined} value={password} onChange={handleOnChange} name="password" className="auth-form__field" label="Пароль" type="password" placeholder="Введите пароль" />
            <Input uiSize="big" error={showErrors ? repeatError! : undefined} value={repeat} onChange={handleOnChange} name="repeat" className="auth-form__field" label="Повторение пароля" type="password" placeholder="Повторите пароль" />

            <Button uiSize="big" disabled={(nameError || passwordError || repeatError) && showErrors ? true : false} className="auth-form__button">Создать</Button>

            <p>Уже зарегистрированы? <span onClick={animateForm}>Войти</span></p>
        </form>
    )
}