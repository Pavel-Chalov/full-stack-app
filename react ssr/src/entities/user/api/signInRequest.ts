import { AuthAPI } from "../../../shared/api/authAPI"
import { TypeAuthData } from "../models/authData"

export const signInRequest = (props: TypeAuthData) => {
    return AuthAPI.post("/sign-in", props)
}