// @flow
import { useLocation } from 'react-router-dom';

const useQuery = (): URLSearchParams => new URLSearchParams(useLocation().search);

export default useQuery;
