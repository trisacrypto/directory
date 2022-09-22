import { Box, BoxProps } from '@chakra-ui/react';
import TrisaLogo from 'assets/trisa_logo.svg';
import CkLazyLoadImage from 'components/LazyImage';
import { Link } from 'react-router-dom';

export default function Logo(props: BoxProps) {
  return (
    <Box {...props}>
      <Link to="/">
        <CkLazyLoadImage src={TrisaLogo} alt="Trisa logo" width="100px" height="50px" />
      </Link>
    </Box>
  );
}
