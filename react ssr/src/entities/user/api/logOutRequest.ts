import { AuthAPI } from "../../../shared/api/authAPI"

export const logOutRequest = () => {
    return AuthAPI.post("/log-out")
}