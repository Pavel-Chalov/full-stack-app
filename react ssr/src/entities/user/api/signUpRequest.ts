import { AuthAPI } from "../../../shared/api/authAPI"
import { TypeAuthData } from "../models/authData"

export const signUpRequest = (props: TypeAuthData) => {
    return AuthAPI.post("/sign-up", props)
}