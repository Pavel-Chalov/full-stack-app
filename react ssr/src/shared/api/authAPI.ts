import axios from "axios";
import serverURL from "../../app/config/serverURL";

export const AuthAPI = axios.create({
    baseURL: `${serverURL}/auth`,
    withCredentials: true
});