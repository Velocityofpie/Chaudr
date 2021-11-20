import{configureStore} from'@reduxjs/toolkit'
import userReducer  from '../components/User/userSlice'
import appReducer from "../components/User/appSlice"

export default configureStore({
    reducer:{
        user: userReducer,
        app: appReducer,
    }
})