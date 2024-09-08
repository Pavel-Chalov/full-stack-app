import { AuthAPI } from "../../../shared/api/authAPI"
import { TypeAuthData } from "../models/authData"

export const changeUserDataRequest = (data: TypeAuthData) => {
    return AuthAPI.put("/change-data", data)
}