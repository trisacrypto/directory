import React from 'react'
import { Modal as RbModal } from 'react-bootstrap'

export const ModalContext = React.createContext(null)

const Modal = (props) => {
    const [isOpen, setIsOpen] = React.useState(false)

    return <ModalContext.Provider value={[isOpen, setIsOpen]} {...props} />
}

const ModalOpenButton = ({ children: child }) => {
    const [, setIsOpen] = React.useContext(ModalContext)
    return React.cloneElement(child, {
        onClick: () => setIsOpen(true)
    })
}

const ModalCloseButton = ({ children: child }) => {
    const [, setIsOpen] = React.useContext(ModalContext)
    return React.cloneElement(child, {
        onClick: () => setIsOpen(false)
    })
}

const ModalContentBase = ({ isOpen, title, children, ...props }) => {

    return (
        <RbModal
            size="lg"
            show={isOpen}
            aria-labelledby="contained-modal-title-vcenter"
            centered
            {...props}
        >
            {children}
        </RbModal>
    )
}

const ModalContent = ({ ...props }) => {
    const [isOpen, setIsOpen] = React.useContext(ModalContext)

    return <ModalContentBase isOpen={isOpen} onHide={() => setIsOpen(false)} {...props} />
}


export { ModalContent, ModalOpenButton, Modal, ModalCloseButton }