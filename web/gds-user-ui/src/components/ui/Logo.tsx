import { Box, BoxProps } from '@chakra-ui/react';
import CkLazyLoadImage from 'components/LazyImage';
import { Link } from 'react-router-dom';
import TrisaLogo from 'assets/TRISA-GDS-black.png';

export default function Logo(props: BoxProps) {
  return (
    <Box {...props}>
      <Link to="/">
        <CkLazyLoadImage src={TrisaLogo} alt="Trisa logo" objectFit="none" />
      </Link>
    </Box>
  );
}
