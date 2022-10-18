import { Box, BoxProps } from '@chakra-ui/react';
import CkLazyLoadImage from 'components/LazyImage';
import { Link } from 'react-router-dom';
import TrisaLogo from 'assets/TRISA-GDS.svg';

export default function Logo(props: BoxProps) {
  return (
    <Box {...props}>
      <Link to="/">
        <CkLazyLoadImage src={TrisaLogo} alt="Trisa logo" objectFit={{ base: 'none' }} />
      </Link>
    </Box>
  );
}
