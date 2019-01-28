<?php
/**
 * 入口
 * Created by PhpStorm.
 * User: 倒霉狐狸
 * Date: 2018-08-06
 * Time: 17:12
 */

error_reporting(E_ALL|E_STRICT);
ini_set('display_errors', 1);
header("Access-Control-Allow-Origin: *");
header("Access-Control-Allow-Methods: *");
header("Access-Control-Allow-Headers: X-Requested-With,LOGIN-KEY,ZHIMIAO-USER,login-key,zhimiao-user");
header("Content-type: text/html; charset=utf-8");

if (isset($_SERVER['REQUEST_METHOD']) && strtolower($_SERVER['REQUEST_METHOD']) == 'options') {
    exit;
}

// 定义纸喵php框架根目录
define('ROOT_PATH', dirname(__DIR__));

require __DIR__.'/../vendor/autoload.php';

new zhimiao\Run;