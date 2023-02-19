import sh


def run(cmd, *args):
    sh.run(cmd, args=args)
