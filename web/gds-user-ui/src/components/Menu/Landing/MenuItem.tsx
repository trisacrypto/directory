import { Link, Text } from '@chakra-ui/react';
import { forwardRef } from 'react';

interface MenuItemProps {
  children: React.ReactNode;
  isLast?: boolean;
  to: string;
}

const MenuItem = forwardRef<any, MenuItemProps>(({ children, isLast, to = '/', ...rest }, ref) => (
  <Text
    mb={{ base: isLast ? 0 : 4, sm: 0 }}
    mr={{ base: 2, sm: isLast ? 8 : 2 }}
    pl={isLast ? 8 : 0}
    display="block"
    as="div"
    {...rest}
    ref={ref}>
    <Link
      isExternal={!!to.startsWith('http')}
      href={to}
      _active={{ outline: 'none' }}
      _focus={{ outline: 'none' }}>
      {children}
    </Link>
  </Text>
));

MenuItem.displayName = 'MenuItem';

export default MenuItem;
