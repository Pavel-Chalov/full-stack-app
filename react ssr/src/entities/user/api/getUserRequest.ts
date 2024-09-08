import { AuthAPI } from "../../../shared/api/authAPI"
import { UserModel } from "../models/userModel"

type Response = {
    message: string,
    user: UserModel
}

export const getUserRequest = async (): Promise<UserModel> => {
    const res = await AuthAPI.get<Response>("/get-data")
    
    return res.data.user
}