import { Meta, Story } from "@storybook/react";
import StatusCard from ".";

interface StatusCardProps {
  mainnetstatus: string;
  testnetstatus: string;
}
export default {
  title: "components/StatusCard",
  component: StatusCard,
} as Meta;

const Template: Story = (args: StatusCardProps) => <StatusCard {...args} />;

export const Default = Template.bind({});
Default.args = {};
