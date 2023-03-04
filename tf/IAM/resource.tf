resource "aws_iam_role" "myecs_task_execution_role" {
  name               = join("-", [var.base_name, "execution-role"])
  assume_role_policy = data.aws_iam_policy_document.myecs_task_execution_assume_policy.json
}

resource "aws_iam_role_policy_attachment" "myecs_task_execution_policy" {
  role       = aws_iam_role.myecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role" "myecs_task_role" {
  name               = join("-", [var.base_name, "role"])
  assume_role_policy = data.aws_iam_policy_document.myecs_task_assume_policy.json
}

resource "aws_iam_policy" "myecs_task_policy" {
  name   = join("-", [var.base_name, "policy"])
  policy = data.aws_iam_policy_document.myecs_task_policy.json
}

resource "aws_iam_role_policy_attachment" "myecs_task_role" {
  role       = aws_iam_role.myecs_task_role.name
  policy_arn = aws_iam_policy.myecs_task_policy.arn
}
