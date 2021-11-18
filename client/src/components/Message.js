import { Avatar } from '@mui/material'
import React from 'react'
import "./Message.css"

function Message() {
    return (
        <div className="message">
            <Avatar/>
            <div className ="message__tag">
                <h4>
                    username
                    {/* <span className="message__timestamp"> timestap</span> */}
                </h4>
                <p>
                    Message goes here
                </p>
            </div>

        </div>
    )
}

export default Message
