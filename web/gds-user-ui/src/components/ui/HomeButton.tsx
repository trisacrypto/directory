import { Button, Icon, Tooltip, ButtonProps } from '@chakra-ui/react';
import { AiOutlineHome } from 'react-icons/ai';

type HomeButtonProps = {
  link: string;
};

const HomeButton: React.FC<HomeButtonProps & ButtonProps> = ({ link, ...props }) => {
  return (
    <Button
      as={'a'}
      href={link}
      {...props}
      backgroundColor={'transparent'}
      _hover={{ background: 'black', color: 'white' }}
      _focus={{ background: 'red.100', color: 'red.500' }}
      borderRadius={0}>
      <Icon as={AiOutlineHome} />
    </Button>
  );
};

export default HomeButton;
