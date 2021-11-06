import styles from "./Counter.css";
import counterReducer from "../components/counterSlice"

export default configureStore({
    reducer:{
        counter: counterReducer,
    }
})