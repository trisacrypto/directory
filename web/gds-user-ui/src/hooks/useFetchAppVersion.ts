
import {
    getAppGitVersion,
    getAppVersionNumber,
    getBffAndGdsVersion
    // isProdEnv
} from 'application/config';
import { useEffect, useState } from 'react';

const useFetchAppVersion = () => {
    const [appVersion, setAppVersion] = useState('');
    const [appGitVersion, setAppGitVersion] = useState('');
    const [bffAndGdsVersion, setBffAndGdsVersion] = useState('');

    useEffect(() => {
        const fetchAppVersion = async () => {
            const appVersionNumber = await getAppVersionNumber() as string;
            const appGitVersionNumber = await getAppGitVersion() as string;
            const bffAndGdsVersionNumber = await getBffAndGdsVersion() as any;

            setAppVersion(appVersionNumber);
            setAppGitVersion(appGitVersionNumber);
            setBffAndGdsVersion(bffAndGdsVersionNumber?.version);
        };

        fetchAppVersion();
    }, []);

    return { appVersion, appGitVersion, bffAndGdsVersion };
};

export default useFetchAppVersion;
