import Axios from './https'

export function Login(data:{[key: string]: any}) {
     return Axios.post('/mock/api/Login',data)
}