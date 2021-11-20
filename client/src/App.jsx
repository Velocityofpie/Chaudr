import React from 'react'
//import {useSelector} from 'react-redux'
import "./App.css"
import Chat from './components/Chat'
import Sidebar from './components/Sidebar/Sidebar'
import RightSidebar from './components/Sidebar/RightSidebar'
//import {selectUser} from "./components/User/userSlice"

const App = () => {
    
    return (
        <div className="app">
            
                <>
                    <Sidebar/>
                    <Chat/>
                    <RightSidebar/>
                </>
            
      </div>
        
    )
}

export default App
