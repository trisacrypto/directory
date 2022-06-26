import { Link, Text } from '@chakra-ui/react';
import { Link as RouteLink } from 'react-router-dom';
interface MenuItemProps {
  children: React.ReactNode;
  isLast?: boolean;
  to?: string;
}

const MenuItem = ({ children, isLast, to = '/', ...rest }: MenuItemProps): JSX.Element => {
  return (
    <Text
      mb={{ base: isLast ? 0 : 4, sm: 0 }}
      mr={{ base: 2, sm: isLast ? 8 : 2 }}
      px={2}
      pl={isLast ? 5 : 0}
      display="block"
      as="div"
      {...rest}>
      {to.startsWith('http') ? (
        <Link href={to} isExternal _active={{ outline: 'none' }} _focus={{ outline: 'none' }}>
          {children}
        </Link>
      ) : (
        <RouteLink to={to}>
          <Link
            isExternal={!!to.startsWith('http')}
            _active={{ outline: 'none' }}
            _focus={{ outline: 'none' }}>
            {children}
          </Link>
        </RouteLink>
      )}
    </Text>
  );
};

export default MenuItem;
