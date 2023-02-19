import sh


def token(*args):
    for arg in args:
        print "%s -> %s" % (arg, sh.tokenise(arg))
