data "aws_vpc" "myecs" {
  // ALBを動かしたいVPCを参照できるように適切に引数を設定してください
}

data "aws_subnets" "myecs_public" {
  // ALBを動かしたいパブリックサブネットを参照できるように適切に引数を設定してください
}

data "aws_subnets" "myecs_private" {
  // ALBを動かしたいプライベートサブネットを参照できるように適切に引数を設定してください
}

data "aws_acm_certificate" "myecs" {
  domain = var.domain_name
}
