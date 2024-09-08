import { SignInForm, SignUpForm } from "../../../features/authForms"
import "./auth.scss"

export const Auth = () => {
    return (
        <div className="auth page centered">
            <div className="auth__container">
                <SignInForm />
                <SignUpForm />
            </div>
        </div>
    )
}