import { DeleteIcon } from '@chakra-ui/icons';
import { Button, Icon, Tooltip, TooltipProps } from '@chakra-ui/react';
import { t } from '@lingui/macro';

type DeleteButtonProps = {
  tooltip?: Omit<TooltipProps, 'children'>;
  onDelete?: (e: unknown) => void;
  isDisabled?: boolean;
};

const TOOLTIPS_DELAY = 2 * 1000;

const DeleteButton: React.FC<DeleteButtonProps> = ({ tooltip, onDelete, isDisabled }) => {
  return (
    <Tooltip label={t`Delete`} openDelay={TOOLTIPS_DELAY} {...tooltip}>
      <Button
        onClick={onDelete}
        visibility={isDisabled ? 'hidden' : 'visible'}
        variant="ghost"
        _hover={{ background: 'red.100', color: 'red.500' }}
        _focus={{ background: 'red.100', color: 'red.500' }}
        borderRadius={0}>
        <Icon as={DeleteIcon} />
      </Button>
    </Tooltip>
  );
};

export default DeleteButton;
