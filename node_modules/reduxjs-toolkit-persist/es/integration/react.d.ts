import { PureComponent, ReactNode } from 'react';
import type { Persistor } from '../types';
declare type Props = {
    onBeforeLift?: () => void;
    children: ReactNode | ((state: boolean) => ReactNode);
    loading: ReactNode;
    persistor: Persistor;
};
declare type State = {
    bootstrapped: boolean;
};
export declare class PersistGate extends PureComponent<Props, State> {
    static defaultProps: {
        children: null;
        loading: null;
    };
    state: {
        bootstrapped: boolean;
    };
    _unsubscribe?: () => void;
    componentDidMount(): void;
    handlePersistorState: () => void;
    componentWillUnmount(): void;
    render(): ReactNode;
}
export {};
