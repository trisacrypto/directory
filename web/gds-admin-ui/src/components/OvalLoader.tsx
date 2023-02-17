import { ReactNode } from 'react';
import Oval from './Oval';

type OvalLoaderProps = {
    title?: ReactNode;
};

function OvalLoader({ title = '', ...rest }: OvalLoaderProps) {
    return (
        <div className="text-center flex flex-column justify-content-center" {...rest}>
            <div>
                <Oval width="35" height="35" stroke="#6b7280" />
            </div>
            <div className="mt-1">
                <small className="block">{title || 'Loading...'}</small>
            </div>
        </div>
    );
}

export default OvalLoader;
