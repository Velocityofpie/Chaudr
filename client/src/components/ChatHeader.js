import React,{useState}  from 'react'
import './ChatHeader.css'
import ExpandMoreIcon from "@material-ui/icons/ExpandMore"
import AddIcon from "@material-ui/icons/Add";


function ChatHeader() {
    return (
        <div className='chatHeader'>
            <div className='chatHeader__left'>
                <h3>
                    <span className='chatHeader_name'> 
                        -
                    </span>
                        test Channel Name
                </h3>
            </div>
            <div className='chatHeader__right'>
                <h3>
                    Add User
                </h3>
                <Navbar className='navbar__right'>
                    <NavItem icon={<AddIcon/>}>
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
                <div className= "Change_ChatName">
                    <form>
                        <input placeholder={'Enter New Change Room Name'}/>
                        <button className="chat__inputButton" type ="submit">
                            send message
                        </button>
                    </form>
                </div>
            </DropdownItem>
            <DropdownItem>
                <div className= "Change_UserName">
                    <form>
                        <input placeholder={'Enter New Username'}/>
                        <button className="chat__inputButton" type ="submit">
                            send message
                        </button>
                    </form>
                </div>
            </DropdownItem>
            <DropdownItem>
                <div className= "Add_User">
                    <form>
                        <input placeholder={'Add New User'}/>
                        <button className="chat__inputButton" type ="submit">
                            send message
                        </button>
                    </form>
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

