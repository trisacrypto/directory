import { Link, Text } from '@chakra-ui/react';

interface MenuItemProps {
  children: React.ReactNode;
  isLast?: boolean;
  to: string;
}

const MenuItem = ({ children, isLast, to = '/', ...rest }: MenuItemProps): JSX.Element => {
  return (
    <Text
      mb={{ base: isLast ? 0 : 4, sm: 0 }}
      mr={{ base: 2, sm: isLast ? 8 : 2 }}
      pl={isLast ? 8 : 0}
      display="block"
      {...rest}>
      <Link
        isExternal={!!to.startsWith('http')}
        href={to}
        _active={{ outline: 'none' }}
        _focus={{ outline: 'none' }}>
        {children}
      </Link>
    </Text>
  );
};

export default MenuItem;
