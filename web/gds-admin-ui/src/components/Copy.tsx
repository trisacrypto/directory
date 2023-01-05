import React from 'react';
import { copyToClipboard } from '@/utils';
import Button from 'react-bootstrap/Button';

type CopyProps = {
    data: any;
};

function Copy({ data }: CopyProps) {
    const [copied, setCopied] = React.useState(false);

    React.useEffect(() => {
        const timeoutId = setTimeout(() => {
            setCopied(false);
        }, 2000);

        return () => {
            clearTimeout(timeoutId);
        };
    });

    const handleCopy = async () => {
        await copyToClipboard(data);
        setCopied(true);
    };
    return copied ? (
        <Button variant="light" className="btn-circle ms-1" title="Copied" data-testid="check-icon">
            <i className="mdi mdi-check text-success" />
        </Button>
    ) : (
        <Button onClick={handleCopy} variant="light" className="btn-circle ms-1" title="Copy" data-testid="copy-icon">
            <i className="mdi mdi-content-copy text-success" />
        </Button>
    );
}

export default Copy;
