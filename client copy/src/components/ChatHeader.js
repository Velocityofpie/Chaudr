import React,{useState}  from 'react'
import './ChatHeader.css'
import PersonAddIcon from '@mui/icons-material/PersonAdd';
import Timer from "./Timer";

function ChatHeader() {
    return (
        <div className='chatHeader'>
            <div className='chatHeader__left'>
                <h3>
                    <span className='chatHeader_name'> 
                        - Channel Name
                    </span>   
                </h3>
            </div>
            <div className='chatHeader__right'>
                <h3>
                    Add User
                </h3>
                <Navbar className='navbar__right'>
                    <NavItem icon={<PersonAddIcon/>}>
                        <DropdownMenu/>
                    </NavItem>
                </Navbar>
            </div>
        </div>
    )
}
function DropdownMenu(){
    function DropdownItem(props){
        return(
            <a href="#" className="menu_item" >
                <span className="icon_button">{props.leftIcon}</span>
                {props.children}
                <span className="icon_right">{props.rightIcon}</span>
            </a>
        )
    }
    return(
        <div className="dropdown">
            <DropdownItem>
                <div className= "Chatroom_Id">
                    Chat room Id:
                </div>
            </DropdownItem>
            <DropdownItem>
                <div className= "Gen_Code">
                    Code:
                </div>
            </DropdownItem>
            <DropdownItem>
                <div className= "Timer">
                    <Timer/>
                </div>
            </DropdownItem>
        </div>
    )
}


function Navbar(props) {
    return (
        <nav className="navbar">
            <ul className="navbar_nav"> {props.children }</ul>
        </nav>
    )
}

function NavItem(props) {
    const[open,setOpen]=useState(false);
    return (
        <li className="nav_item">
            <a href="#" className="icon_button" onClick={() => setOpen(!open)}>
                {props.icon}
            </a>
            {open && props.children}
        </li>
    )
}



export default ChatHeader

