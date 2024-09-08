import { AuthAPI } from "../../../shared/api/authAPI"

export const refreshRequest = async () => {
    try {
        const response = await AuthAPI.get("/refresh");
        return response.data;
      } catch (error) {
        throw error;
      }
}