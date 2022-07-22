import { chakra, ChakraComponent } from '@chakra-ui/react';
import { ForwardRefExoticComponent, RefAttributes } from 'react';
import { Link as RouterLink, LinkProps } from 'react-router-dom';

const ChakraRouterLink: ChakraComponent<
  ForwardRefExoticComponent<LinkProps & RefAttributes<HTMLAnchorElement>>,
  {}
> = chakra(RouterLink);

export default ChakraRouterLink;
