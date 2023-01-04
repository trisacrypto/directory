import React from 'react';

const modalContext = React.createContext();

export const actionType = {
  SEND_EMAIL_MODAL: 'SEND_EMAIL_MODAL',
  CLOSE_MODAL: 'CLOSE_MODAL',
};
const reducer = (state, action) => {
  switch (action.type) {
    case actionType.SEND_EMAIL_MODAL:
      return { ...state, toggle: true, ...action.payload };
    case actionType.CLOSE_MODAL:
      return { ...state, toggle: false };
    default:
      throw new Error(`unhandled type ${action.type}`);
  }
};

const ModalProvider = (props) => {
  const [state, dispatch] = React.useReducer(reducer, {
    toggle: false,
    vasp: { name: '', id: '' },
  });

  const value = {
    ...state,
    dispatch,
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
