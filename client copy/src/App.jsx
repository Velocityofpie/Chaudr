import React from 'react'
import "./App.css"
import Chat from './components/Chat'
import Sidebar from './components/Sidebar/Sidebar'
import RightSidebar from './components/Sidebar/RightSidebar'

const App = () => {
    return (
        <div className="app">
            <Sidebar/>
            <Chat/>
            <RightSidebar/>

      </div>
        
    )
}

export default App
