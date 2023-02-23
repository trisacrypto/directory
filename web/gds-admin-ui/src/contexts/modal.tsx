import React from 'react';

type Payload = { name: string; id: string };
type State = {
    toggle: boolean;
    vasp: Payload;
    openSendEmailModal: (vasp: Payload) => void;
    closeSendEmailModal: () => void;
};
const modalContext = React.createContext<State | null>(null);

export const actionType = {
    SEND_EMAIL_MODAL: 'SEND_EMAIL_MODAL',
    CLOSE_MODAL: 'CLOSE_MODAL',
};

type Action = {
    type: keyof typeof actionType;
    payload?: any;
};

const reducer = (state: State, action: Action) => {
    switch (action.type) {
        case actionType.SEND_EMAIL_MODAL:
            return { ...state, toggle: true, ...action.payload };
        case actionType.CLOSE_MODAL:
            return { ...state, toggle: false };
        default:
            throw new Error(`unhandled type ${action.type}`);
    }
};

const ModalProvider = (props: any) => {
    const [state, dispatch] = React.useReducer(reducer, {
        toggle: false,
        vasp: { name: '', id: '' },
    });

    const openSendEmailModal = (vasp: Payload) => dispatch({ type: 'SEND_EMAIL_MODAL', payload: { vasp } });

    const closeSendEmailModal = () => dispatch({ type: 'CLOSE_MODAL' });

    const value = {
        ...state,
        openSendEmailModal,
        closeSendEmailModal,
    };

    return <modalContext.Provider value={value} {...props} />;
};

const useModal = () => {
    const context = React.useContext(modalContext);
    if (!context) {
        throw new Error('useModal should be used within a ModalProvider');
    }

    return context;
};

export { ModalProvider, useModal };
