import React, { ReactNode, useContext } from 'react';
import { Modal as RbModal } from 'react-bootstrap';

type Value = {
    openModal: () => void;
    closeModal: () => void;
    isOpen: boolean;
};

type ModalContentBaseProps = {
    title?: ReactNode;
    isOpen?: boolean;
    children?: ReactNode;
    onHide?: () => void;
};

export const ModalContext = React.createContext<Value | null>(null);

export const Modal = (props: any) => {
    const [isOpen, setIsOpen] = React.useState(false);

    const value = {
        openModal: () => setIsOpen(true),
        closeModal: () => setIsOpen(false),
        isOpen,
    };

    return <ModalContext.Provider value={value} {...props} />;
};

export const ModalOpenButton = ({ children: child }: any) => {
    const { openModal } = useModalContext();
    return React.cloneElement(child, {
        onClick: openModal,
    });
};

export const ModalCloseButton = ({ children: child }: any) => {
    const { closeModal } = useModalContext();

    return React.cloneElement(child, {
        onClick: closeModal,
    });
};

export const ModalContentBase = ({ isOpen, title, children, ...props }: ModalContentBaseProps) => {
    return (
        <RbModal size="lg" show={isOpen} aria-labelledby="contained-modal-title-vcenter" centered {...props}>
            {children}
        </RbModal>
    );
};

export const ModalContent = ({ ...props }) => {
    const { closeModal, isOpen } = useModalContext();

    if (!isOpen) return null;

    return <ModalContentBase isOpen={isOpen} onHide={closeModal} {...props} />;
};

export function useModalContext() {
    const context = useContext(ModalContext);

    if (!context) {
        throw new Error(`useModalContext should be called within a ModalContext`);
    }

    return context;
}
