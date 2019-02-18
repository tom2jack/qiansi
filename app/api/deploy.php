<?php
/**
 * 部署服务模块
 * Created by PhpStorm.
 * User: 倒霉狐狸
 * Date: 2019-02-18
 * Time: 19:31
 */

namespace app\api;


class deploy
{
    private $uid;
    public function __construct()
    {
        // 获取用户UID并同时判断登陆
        $this->uid = \app\Service\Verify::isLogin();
    }

    public function index()
    {
        
    }
}